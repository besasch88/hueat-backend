package order

import (
	"net"

	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_pubsub"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/justinmichaelvieira/escpos"
	"gorm.io/gorm"
)

type printServiceInterface interface {
	print(ctx *gin.Context, input printOrderInputDto) error
}

type printService struct {
	printerEnabled  bool
	storage         *gorm.DB
	pubSubAgent     *ceng_pubsub.PubSubAgent
	repository      orderRepositoryInterface
	printRepository printRepositoryInterface
}

func newPrintService(printerEnabled bool, storage *gorm.DB, pubSubAgent *ceng_pubsub.PubSubAgent, repository orderRepositoryInterface, printRepository printRepositoryInterface) printServiceInterface {
	return printService{
		printerEnabled:  printerEnabled,
		storage:         storage,
		pubSubAgent:     pubSubAgent,
		repository:      repository,
		printRepository: printRepository,
	}
}

func (s printService) print(ctx *gin.Context, input printOrderInputDto) error {
	switch input.Target {
	case "order":
		return s.printOrder(uuid.MustParse(input.TableID))
	case "course":
		return s.printCourse(uuid.MustParse(input.TableID), uuid.MustParse(*input.CourseID))
	case "bill":
		return s.printBill(uuid.MustParse(input.TableID))
	case "payment":
		return s.printPayment(uuid.MustParse(input.TableID))

	default:
		return errInvalidPrintRequest
	}
}

func (s printService) printOrder(tableId uuid.UUID) error {
	items, err := s.repository.getOrderDetailByTableID(s.storage, tableId)
	if err != nil {
		return err
	}
	if !s.printerEnabled {
		return nil
	}
	return s.printItems(items)
}

func (s printService) printCourse(tableId uuid.UUID, courseId uuid.UUID) error {
	items, err := s.repository.getOrderDetailByTableIDAndCourseID(s.storage, tableId, courseId)
	if err != nil {
		return err
	}
	if !s.printerEnabled {
		return nil
	}
	return s.printItems(items)
}

func (s printService) printBill(tableId uuid.UUID) error {
	items, err := s.repository.getPricedOrderByTableID(s.storage, tableId)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	if !s.printerEnabled {
		return nil
	}
	conn, err := net.Dial("tcp", items[0].PrinterURL)
	if err != nil {
		return err
	}
	defer conn.Close()
	printer := escpos.New(conn)
	if err := s.printRepository.printTable(printer, items[0].TableName); err != nil {
		return err
	}
	if err := s.printRepository.printPrinterName(printer, items[0].PrinterTitle); err != nil {
		return err
	}
	if err := s.printRepository.printTableCreation(printer, items[0].Username, items[0].TableCreatedAt); err != nil {
		return err
	}
	if err := s.printRepository.printLine(printer); err != nil {
		return err
	}
	total := int64(0)
	for _, item := range items {
		if item.MenuOptionTitle != nil {
			if err := s.printRepository.printItemAndPrice(printer, item.Quantity, *item.MenuOptionTitle, *item.MenuOptionPrice); err != nil {
				return err
			}
			total = total + (item.Quantity * *item.MenuOptionPrice)
		} else {
			if err := s.printRepository.printItemAndPrice(printer, item.Quantity, item.MenuItemTitle, item.MenuItemPrice); err != nil {
				return err
			}
			total = total + (item.Quantity * item.MenuItemPrice)
		}
	}
	if err := s.printRepository.printLine(printer); err != nil {
		return err
	}
	if err := s.printRepository.printTotalPrice(printer, total); err != nil {
		return err
	}
	if err := s.printRepository.printLine(printer); err != nil {
		return err
	}
	if err := s.printRepository.printRecipeCollection(printer); err != nil {
		return err
	}
	if err := s.printRepository.printAndCut(printer); err != nil {
		return err
	}
	return nil
}

func (s printService) printPayment(tableId uuid.UUID) error {
	item, err := s.repository.getTotalPriceAndPaymentByTableID(s.storage, tableId)
	if err != nil {
		return err
	}
	if ceng_utils.IsEmpty(item) {
		return nil
	}
	if !s.printerEnabled {
		return nil
	}
	conn, err := net.Dial("tcp", item.PrinterURL)
	if err != nil {
		return err
	}
	defer conn.Close()
	printer := escpos.New(conn)
	if err := s.printRepository.printTable(printer, item.TableName); err != nil {
		return err
	}
	if err := s.printRepository.printPrinterName(printer, item.PrinterTitle); err != nil {
		return err
	}
	if err := s.printRepository.printTableCreation(printer, item.Username, item.TableCreatedAt); err != nil {
		return err
	}
	if err := s.printRepository.printLine(printer); err != nil {
		return err
	}
	if err := s.printRepository.printPaymentMethod(printer, item.TablePayment, item.PriceTotal); err != nil {
		return err
	}
	if err := s.printRepository.printAndCut(printer); err != nil {
		return err
	}
	return nil
}

func (s printService) printItems(items []orderDetailEntity) error {
	var conn net.Conn
	var err error
	var printer *escpos.Escpos
	lastPrinterTitle := ""
	lastCourseID := ""
	for _, item := range items {
		if item.PrinterTitle != lastPrinterTitle {
			// if the printer is changing, send the print and cut and close the previous connection
			if printer != nil {
				err = s.printRepository.printAndCut(printer)
				if err != nil {
					return err
				}
			}
			if conn != nil {
				conn.Close()
			}
			// So create a connection to the new printer
			conn, err = net.Dial("tcp", item.PrinterURL)
			if err != nil {
				return err
			}
			printer = escpos.New(conn)
			lastPrinterTitle = item.PrinterTitle
			lastCourseID = ""
			if err := s.printRepository.printTable(printer, item.TableName); err != nil {
				return err
			}
			if err := s.printRepository.printPrinterName(printer, item.PrinterTitle); err != nil {
				return err
			}
			if err := s.printRepository.printTableCreation(printer, item.Username, item.TableCreatedAt); err != nil {
				return err
			}
		}
		if item.CourseID != lastCourseID {
			lastCourseID = item.CourseID
			// The course changed, so print it on paper
			if err := s.printRepository.printCourse(printer, item.CourseNumber); err != nil {
				return err
			}
		}
		// Now for each element, I can print them
		if item.MenuOptionTitle != nil {
			if err := s.printRepository.printItem(printer, item.Quantity, *item.MenuOptionTitle); err != nil {
				return err
			}
		} else {
			if err := s.printRepository.printItem(printer, item.Quantity, item.MenuItemTitle); err != nil {
				return err
			}
		}
	}
	// At the end, if needed, print and cut and close the connection
	if conn != nil && printer != nil {
		if err := s.printRepository.printAndCut(printer); err != nil {
			return err
		}
		conn.Close()
	}
	return nil
}

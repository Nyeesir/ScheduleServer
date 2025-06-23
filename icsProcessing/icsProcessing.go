package icsProcessing

//TODO: Udar mózgu chyba, zamienić na mapy
//NIESKOŃCZONA PĘTLA

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrCantReadLine = errors.New("can't read line")
	ErrNotValidFile = errors.New("not a valid file")
	ErrEndOfFile    = errors.New("end of file")
)

func IcsToJson(icsString string) ([]byte, error) {
	//callendar := CalendarTemplate{}
	// scanner := bufio.NewScanner(strings.NewReader(icsString))

	//if !scanner.Scan() {
	//	return nil, ErrCantReadLine
	//}
	//
	//if !(scanner.Text() == "BEGIN:VCALENDAR") {
	//	return nil, ErrNotValidFile
	//}
	//
	//parseMainProperties(scanner, &callendar)
	////parseEvents(scanner, &callendar)
	//event, err := parseNextEvent(scanner)
	//if err != nil {
	//	return nil, err
	//}
	//callendar.Events = append(callendar.Events, event)
	//
	//return json.Marshal(&callendar)
}

//func parseMainProperties(scanner *bufio.Scanner, callendar *CalendarTemplate) {
//	for scanner.Scan() {
//		line := scanner.Text()
//		sides := strings.Split(line, ":")
//		left := sides[0]
//		right := sides[1]
//		switch left {
//		case "PRODID":
//			callendar.ProdId = right
//		case "VERSION":
//			callendar.Version = right
//		case "CALSCALE":
//			callendar.CalScale = right
//		case "METHOD":
//			callendar.Method = right
//		case "X-WR-CALNAME":
//			callendar.CalName = right
//		case "X-WR-TIMEZONE":
//			callendar.CalTimeZone = right
//		default:
//			return
//		}
//	}
//}
//
//func parseNextEvent(scanner *bufio.Scanner) (EventTemplate, error) {
//	event := EventTemplate{}
//	found := false
//	fmt.Println("jestem tutaj")
//	for scanner.Scan() {
//		line := scanner.Text()
//		fmt.Println(line)
//		fmt.Println("tutaj tez")
//		if line == "BEGIN:VEVENT" {
//			found = true
//			break
//		} else if line == "BEGIN:VEVENT" {
//			return event, ErrEndOfFile
//		}
//	}
//
//	if found {
//		for scanner.Scan() {
//			line := scanner.Text()
//			sides := strings.Split(line, ":")
//			left := sides[0]
//			right := sides[1]
//			fmt.Println(left)
//			switch left {
//			case "DTSTART":
//				stampTime, err := time.Parse("T150405Z", right)
//				if err == nil {
//					event.DateTimeStart = stampTime.Format(time.RFC3339)
//				} else {
//					event.DateTimeStamp = ""
//				}
//			case "DTEND":
//				stampTime, err := time.Parse("T150405Z", right)
//				if err == nil {
//					event.DateTimeStart = stampTime.Format(time.RFC3339)
//				} else {
//					event.DateTimeStamp = ""
//				}
//			case "DTSTAMP":
//				stampTime, err := time.Parse("20060102T150405Z", right)
//				if err == nil {
//					event.DateTimeStart = stampTime.Format(time.RFC3339)
//				} else {
//					event.DateTimeStamp = ""
//				}
//			case "UID":
//				event.Uid = right
//			case "CLASS":
//				event.Classification = right
//			case "SEQUENCE":
//				event.Sequence = right
//			case "STATUS":
//				event.Status = right
//			case "SUMMARY":
//				event.Summary = right
//			case "TRANSP":
//				event.Transparency = right
//			case "END":
//				break
//			}
//		}
//	}
//
//	return event, nil
//}
//
//func parseEvents(scanner *bufio.Scanner, callendar *CalendarTemplate) {
//	for true {
//		event, err := parseNextEvent(scanner)
//		if err != nil {
//			break
//		}
//		callendar.Events = append(callendar.Events, event)
//	}
//}

package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func main() {

	buf := bytes.NewBuffer(nil)
	//buf.Write(ESC)
	//buf.Write(GS)
	//
	//buf.Write(ESC)
	//buf.Write([]byte("d"))
	//buf.Write([]byte(strconv.Itoa(10)))
	//
	//buf.Write(GS)
	//buf.Write([]byte("V"))
	//buf.Write(CUT_FULL)
	//buf.Write([]byte(strconv.Itoa(3)))

	//const string ESC = "\u001B"
	//const string GS = "\u001D"
	//const string InitializePrinter = ESC + "@"
	//const string BoldOn = ESC + "E" + "\u0001"
	//const string BoldOff = ESC + "E" + "\0"
	//const string DoubleOn = GS + "!" + "\u0011" // 2x sized text (double-high + double-wide)
	//const string DoubleOff = GS + "!" + "\0"

	//buf.Write(InitializePrinter)
	//buf.WriteString("Here is some normal text.")

	//printJob.Print(InitializePrinter)
	//printJob.PrintLine("Here is some normal text.")
	//printJob.PrintLine(BoldOn + "Here is some bold text." + BoldOff)
	//printJob.PrintLine(DoubleOn + "Here is some large text." + DoubleOff)
	//
	//printJob.ExecuteAsync()

	//* Initialize */
	//$printer -> initialize();
	//
	///* Text */
	//$printer -> text("Hello world\n");
	//$printer -> cut();

	// tspl works
	//buf.WriteString("SIZE 42mm, 42mm, 0 DIRECTION 1\r\n")
	//buf.WriteString("CLS\r\n")
	//buf.WriteString("TEXT 10, 10, \"2\",0,1,1 \"Code 128, switch code\"\r\n")
	//buf.WriteString("PRINT 1,1\r\n")

	//buf.Write(InitializePrinter)
	//
	//buf.Write([]byte("kekw\n\n"))
	//
	//buf.Write(PrintAndFeed)
	//buf.Write([]byte("25"))

	//buf.Write([]byte("\x1Bia\x00"))

	buf.Write([]byte("\x1d\x6b\x041234\x00"))

	println(hex.Dump(buf.Bytes()))
	bts := buf.Bytes()
	parts := len(bts) / 20
	println("parts and len ", parts, len(bts))

	//if err := print(bts); err != nil {
	//	log.Println("err print: ", err)
	//}
	//return

	printerUUID, err := bluetooth.ParseUUID("280f43ba-b547-3d72-c325-a8807b20536c")
	must("failed to parse printer uuid: %w", err)
	serviceUUID, err := bluetooth.ParseUUID("E7810A71-73AE-499D-8C15-FAA9AEF0C3F2")
	must("service uuid", err)
	characteristicUUID, err := bluetooth.ParseUUID("bef8d6c9-9c21-4c9e-b632-bd58c1009f9f")
	must("characteristic uuid", err)

	// Enable BLE interface.
	must("enable BLE stack", adapter.Enable())
	dev, err := adapter.Connect(bluetooth.Address{UUID: printerUUID}, bluetooth.ConnectionParams{
		ConnectionTimeout: bluetooth.NewDuration(time.Second * 30),
		MinInterval:       bluetooth.NewDuration(time.Second * 1),
		MaxInterval:       bluetooth.NewDuration(time.Second * 15),
	})
	must("connect", err)
	println("connected")
	defer func() {
		must("disconnect", dev.Disconnect())
	}()

	srvs, err := dev.DiscoverServices([]bluetooth.UUID{serviceUUID})
	must("failed to discover services:", err)
	if len(srvs) != 1 {
		log.Println("len srvs not 1", srvs)
		return
	}

	srv := srvs[0]

	var chars []bluetooth.DeviceCharacteristic
	chars, err = srv.DiscoverCharacteristics([]bluetooth.UUID{characteristicUUID})
	must("failed to discover characteristics:", err)
	if len(chars) != 1 {
		log.Println("len chars not 1", srvs)
		return
	}

	//buf := make([]byte, 255)
	char := chars[0]

	char.WriteWithoutResponse(bts)

	//if err = char.EnableNotifications(func(buf []byte) {
	//	log.Println("notify: ", hex.Dump(buf))
	//}); err != nil {
	//	log.Println("error to enable notifications:", err)
	//	return
	//}
	//prev := 0
	//for i := 0; i <= parts; i++ {
	//	part := bts[prev : prev+20]
	//	//part = bytes.Trim(part, "\x00")
	//	prev += 20
	//	println(hex.Dump(part))
	//
	//	if _, err = char.WriteWithoutResponse(bts); err != nil {
	//		log.Println("failed to write", err)
	//		return
	//	}
	//
	//	time.Sleep(time.Millisecond * 15)
	//}
	time.Sleep(time.Second * 5)
}

func print(data []byte) error {
	conn, err := net.Dial("tcp", "localhost:9100")
	if err != nil {
		return err
	}

	w, err := conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	log.Println("written: ", w)

	return nil
}

var (
	SELFTEST = []byte("SELFTEST\r\n")
	LF       = []byte{0x0a}
	ESC      = []byte{0x1b}
	NUL      = []byte{0x00}
	GS       = []byte{0x1d}
	CUT_FULL = []byte{0x41}

	InitializePrinter = append(ESC, '@') //ESC + "@"
	PrintAndFeed      = append(ESC, 'j') //ESC + "@"
)

// The UUID for the white lable BLE printer, obtained using LighBlue app.
// The same may be obtained using any other Bluetooth explorer app.
//var serviceUUID = NSUUID(UUIDString: "E7810A71-73AE-499D-8C15-FAA9AEF0C3F2")
//var characteristicUUID = NSUUID(UUIDString: "BEF8D6C9-9C21-4C9E-B632-BD58C1009F9F")

// UART Microchip RN_BLE RN4871
//------- service  49535343-fe7d-4ae5-8fa9-9fafd205e455
//-- 49535343-1e4d-4bd9-ba61-23c647249616    err: Reading is not permitted.
//-- 49535343-8841-43f4-a8d4-ecbe34729bb3    err: Reading is not permitted.

//found device: 280f43ba-b547-3d72-c325-a8807b20536c -47 XP-365B
//connected
//-- characteristic 49535343-1e4d-4bd9-ba61-23c647249616
//Reading is not permitted.
//-- characteristic 49535343-8841-43f4-a8d4-ecbe34729bb3
//Reading is not permitted.
//-- characteristic 00002af0-0000-1000-8000-00805f9b34fb (NOTIFY)
//Reading is not permitted.
//-- characteristic 00002af1-0000-1000-8000-00805f9b34fb (READ, WRITE)
//Reading is not permitted.
//-- characteristic bef8d6c9-9c21-4c9e-b632-bd58c1009f9f (INDICATE, NOTIFY, READ, WRITE , WRITE NO RESPONSE)
//value =

// bef8d6c9-9c21-4c9e-b632-bd58c1009f9f

// chunks split the slice `s` in chunks of the given size.
func chunks(s []byte, size int) [][]byte {
	var result [][]byte
	l := len(s)

	for i := 0; i < l; i += size {
		end := i + size
		if end > l {
			end = l
		}
		result = append(result, s[i:end])
	}

	return result
}

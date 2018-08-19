package messages

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/newrushbolt/AppleSeal/logger"
	chart "github.com/wcharczuk/go-chart"
	"gopkg.in/telegram-bot-api.v4"
)

func ParseMessage(bot *tgbotapi.BotAPI, inMsg *tgbotapi.Message) {
	if inMsg == nil {
		return
	}
	logger.Logger.Debugf("[%s] %s", inMsg.From.UserName, inMsg.Text)
	if inMsg.IsCommand() {
		msgCmd := inMsg.Command()
		msgCmdArgs := inMsg.CommandArguments()
		logger.Logger.Debugf("Got command:\t%s", msgCmd)
		err := error(nil)
		if msgCmd == "test" {
			err = tgCommandTest(bot, inMsg.Chat.ID, msgCmd, msgCmdArgs)
		} else if msgCmd == "temp" {
			err = tgGetTextTemp(bot, inMsg.Chat.ID, msgCmd, msgCmdArgs)
		} else if msgCmd == "temps" {
			err = tgGetTextTemps(bot, inMsg.Chat.ID, msgCmd, msgCmdArgs)
		} else if msgCmd == "image" {
			err = tgGetImageTest(bot, inMsg.Chat.ID, msgCmd, msgCmdArgs)
		} else if msgCmd == "tempimage" {
			err = tgGetGraphTemp(bot, inMsg.Chat.ID, msgCmd, msgCmdArgs)
		}
		if err != nil {
			logger.Logger.Errorf("Cannot send a message:\t%v", err)
		}
	}
}

func getFakeTemp() map[string]float64 {
	rnd1 := rand.Float64() * 50
	rnd2 := rand.Float64() * 50
	res := map[string]float64{
		"Temp1": rnd1,
		"Temp2": rnd2,
	}
	logger.Logger.Debugf("Got fake temp:\n%v", res)
	return res
}

type tempsArray map[string][]float64

func getFakeTempArray(arrDepth int) (tempsArray, []time.Time) {
	data := tempsArray{
		"Temp1": []float64{},
		"Temp2": []float64{},
	}
	timing := []time.Time{}

	for i := 0; i <= arrDepth; i++ {
		data["Temp1"] = append(data["Temp1"], (rand.Float64() * 50))
		data["Temp2"] = append(data["Temp2"], (rand.Float64() * 50))
		timing = append(timing, time.Now().Add((time.Duration(i) * time.Minute * -1)))
	}

	return data, timing
}

func tgGetTextTemp(bot *tgbotapi.BotAPI, chatId int64, cmd string, cmdArgs string) error {
	tempData := getFakeTemp()
	retText := string("")
	for key, value := range tempData {
		retText = retText + fmt.Sprintf("<%s>:\t%.2f\n", key, value)
	}
	msg := tgbotapi.NewMessage(chatId, retText)
	_, err := bot.Send(msg)
	return err
}

func tgGetTextTemps(bot *tgbotapi.BotAPI, chatId int64, cmd string, cmdArgs string) error {
	tempData, _ := getFakeTempArray(10)
	retText := string("")
	for key, value := range tempData {
		retText = retText + fmt.Sprintf("<%s>:\t%v\n", key, value)
	}
	msg := tgbotapi.NewMessage(chatId, retText)
	_, err := bot.Send(msg)
	return err
}

func tgGetGraphTemp(bot *tgbotapi.BotAPI, chatId int64, cmd string, cmdArgs string) error {
	tempData, timing := getFakeTempArray(200)
	logger.Logger.Debugf("Generated fake data:\t%v", tempData)
	logger.Logger.Debugf("Generated timing map:\t%v", timing)

	rndFileName := fmt.Sprintf("/tmp/AppleSealFile-%d.png", rand.Int())
	myfile, err := os.Create(rndFileName)
	if err != nil {
		return err
	}
	defer myfile.Close()

	var myLittleSeries []chart.Series
	for dKey, dValues := range tempData {
		series := chart.TimeSeries{
			Name:    dKey,
			YValues: dValues,
			XValues: timing,
		}
		myLittleSeries = append(myLittleSeries, series)
	}
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style:          chart.StyleShow(),
			TickPosition:   chart.TickPositionBetweenTicks,
			ValueFormatter: chart.TimeMinuteValueFormatter,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  30,
				Left: 30,
			},
		},
		Series: myLittleSeries,
	}
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	err = graph.Render(chart.PNG, myfile)
	if err != nil {
		return err
	}

	logger.Logger.Debugf("Wrote file:\t%s", rndFileName)
	msg := tgbotapi.NewDocumentUpload(chatId, rndFileName)
	_, err = bot.Send(msg)
	return err
}

func tgGetImageTest(bot *tgbotapi.BotAPI, chatId int64, cmd string, cmdArgs string) error {
	d2 := []byte{115, 111, 109, 101, 10}
	rndFileName := fmt.Sprintf("/tmp/AppleSealFile-%d.png", rand.Int())
	err := ioutil.WriteFile(rndFileName, d2, 0644)
	logger.Logger.Debugf("Wrote file:\t%s", rndFileName)
	msg := tgbotapi.NewDocumentUpload(chatId, "/home/svmikhailov/Pictures/apple.jpg")
	_, err = bot.Send(msg)
	return err
}

func tgCommandTest(bot *tgbotapi.BotAPI, chatId int64, cmd string, cmdArgs string) error {
	msg := tgbotapi.NewMessage(chatId, ("Tested:\t" + cmdArgs))
	_, err := bot.Send(msg)
	return err
}

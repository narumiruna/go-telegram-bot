package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/narumiruna/go-visa-fx-rates/pkg/visa"
	log "github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type VISAService struct {
	client *visa.RestClient
}

func NewVISAService() *VISAService {
	return &VISAService{
		client: visa.NewRestClient(),
	}
}

func (s *VISAService) Handle(c tele.Context) error {
	payload := c.Message().Payload

	splits := strings.Split(payload, " ")
	if len(splits) != 3 {
		return c.Reply("Usage: /visa <amount> <from> <to>")
	}

	amount, err := strconv.ParseFloat(splits[0], 64)
	if err != nil {
		return err
	}

	convertedAmount, err := s.rate(amount, splits[1], splits[2])
	if err != nil {
		return err
	}

	return c.Reply(fmt.Sprintf("%f %s = %f %s", amount, splits[1], convertedAmount, splits[2]))
}

func (s *VISAService) rate(amount float64, from string, to string) (float64, error) {
	ctx := context.Background()

	request := visa.RatesRequest{
		Amount:           amount,
		ToCurr:           strings.ToUpper(from),
		FromCurr:         strings.ToUpper(to),
		UTCConvertedDate: time.Now(),
	}

	response, err := s.client.Rates(ctx, request)
	log.Infof("%+v", response)
	if err != nil {
		time.Sleep(1 * time.Second)

		// retry with yesterday's date
		request := visa.RatesRequest{
			Amount:           amount,
			ToCurr:           strings.ToUpper(from),
			FromCurr:         strings.ToUpper(to),
			UTCConvertedDate: time.Now().AddDate(0, 0, -10),
		}

		response, err := s.client.Rates(ctx, request)
		if err != nil {
			log.Infof("%+v", response)
			return 0, err
		}

		return strconv.ParseFloat(response.ConvertedAmount, 64)
	}
	return strconv.ParseFloat(response.ConvertedAmount, 64)
}

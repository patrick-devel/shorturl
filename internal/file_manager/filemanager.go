package filemanager

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	"github.com/patrick-devel/shorturl/internal/models"
)

var ErrNotFoundEvent = errors.New("event not found in file")

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteEvent(event *models.Event) error {
	return p.encoder.Encode(&event)
}

func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *Consumer) ReadEvent(hash string) (*models.Event, error) {
	for c.scanner.Scan() {
		data := c.scanner.Bytes()

		event := models.Event{}
		err := json.Unmarshal(data, &event)
		if err != nil {
			return nil, err
		}

		if event.Hash == hash {
			return &event, nil
		}
	}

	if c.scanner.Err() != nil {
		return nil, c.scanner.Err()
	}

	return nil, ErrNotFoundEvent
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

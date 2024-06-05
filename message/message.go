package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func Messanger(c *fiber.Ctx) error {
	type requestBody struct {
		Text string `json:"text"`
	}

	var reqBody requestBody
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	url := "http://192.168.0.143/message"
	data := reqBody
	jsonData, err := json.Marshal(data)
	println(string(jsonData))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to marshal JSON data",
		})
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to create HTTP request: %s", err.Error()),
		})
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to send data to URL: %s", err.Error()),
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Received non-200 response code from the server: %d", resp.StatusCode),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "message send succesfully",
	})
}

package convert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ConvertText(c *fiber.Ctx) error {
	type requestBody struct {
		Text string `json:"text"`
	}

	var reqBody requestBody
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	morsDict := map[string]string{
		"A": ".-", "B": "-...", "C": "-.-.", "D": "-..", "E": ".", "F": "..-.", "G": "--.", "H": "....", "I": "..", "J": ".---", "K": "-.-", "L": ".-..", "M": "--", "N": "-.", "O": "---", "P": ".--.", "Q": "--.-", "R": ".-.", "S": "...", "T": "-", "U": "..-", "V": "...-", "W": ".--", "X": "-..-", "Y": "-.--", "Z": "--..",
	}

	convertedText := ""
	for _, char := range strings.ToUpper(reqBody.Text) {
		if char == ' ' {
			convertedText += " "
			continue
		}
		if morseCode, exists := morsDict[string(char)]; exists {
			convertedText += morseCode + " "
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid character in text",
			})
		}
	}

	url := "http://192.168.0.143/convert"
	data := map[string]string{"converted_text": strings.TrimSpace(convertedText)}
	jsonData, err := json.Marshal(data)
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
		"converted_text": strings.TrimSpace(convertedText),
	})
}

package services

import (
	"fmt"
	"strings"
)

// GenerateAIResponse generates a mock AI response for now
// In a real application, this would call an external AI service
func GenerateAIResponse(userMessage string) (string, error) {
	// This is a simple mock implementation
	userMessage = strings.ToLower(userMessage)

	if strings.Contains(userMessage, "hello") || strings.Contains(userMessage, "hi") {
		return "Hello! How can I assist you today?", nil
	}

	if strings.Contains(userMessage, "help") {
		return "I'm here to help. What questions do you have?", nil
	}

	if strings.Contains(userMessage, "configure") || strings.Contains(userMessage, "configuration") {
		return "To configure your device, please follow these steps:\n\n1. Connect to the admin panel using the IP address 192.168.1.1\n2. Login with your administrator credentials\n3. Navigate to the 'Settings' tab\n4. Adjust your configuration as needed", nil
	}

	if strings.Contains(userMessage, "update") || strings.Contains(userMessage, "firmware") {
		return "Yes, you can update the firmware remotely. Please follow these steps:\n\n1. Ensure your device is connected to the internet\n2. Access the admin panel\n3. Go to 'System' > 'Updates'\n4. Click 'Check for Updates'\n5. If an update is available, click 'Download and Install'", nil
	}

	if strings.Contains(userMessage, "credentials") || strings.Contains(userMessage, "password") {
		return "Your administrator credentials should have been provided with your device. If you've lost them, you can:\n\n1. Check the documentation that came with your device\n2. Look for a sticker on the device itself\n3. Contact customer support with your device serial number", nil
	}

	// Default response
	return fmt.Sprintf("Thank you for your message: \"%s\". I'm processing your request and will get back to you shortly.", userMessage), nil
}
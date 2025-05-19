package tgbot

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/auth/service"
	userService "github.com/InTeamDev/utmn-map-go-backend/internal/domain/user/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ResponseChannelKey is the key for storing response channels in context
type ResponseChannelKey struct{}

// TelegramBot is the Telegram bot service
type TelegramBot struct {
	bot             *tgbotapi.BotAPI
	authService     *service.AuthService
	userService     *userService.Service
	developersChat  int64
	responseChans   map[string]chan string
	responseChMutex sync.RWMutex
	botUsername     string
}

// Config represents the configuration for the Telegram bot
type Config struct {
	Token          string
	DevelopersChat int64
}

// NewTelegramBot creates a new Telegram bot service
func NewTelegramBot(config Config, authService *service.AuthService, userService *userService.Service) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Telegram bot: %w", err)
	}

	// Get bot information to store the username
	botInfo, err := bot.GetMe()
	if err != nil {
		return nil, fmt.Errorf("failed to get bot info: %w", err)
	}

	return &TelegramBot{
		bot:            bot,
		authService:    authService,
		userService:    userService,
		developersChat: config.DevelopersChat,
		responseChans:  make(map[string]chan string),
		botUsername:    botInfo.UserName,
	}, nil
}

// GetBotUsername returns the bot's username
func (t *TelegramBot) GetBotUsername() string {
	return t.botUsername
}

// Start starts the Telegram bot
func (t *TelegramBot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		go t.handleUpdate(update)
	}
}

// handleUpdate processes incoming updates
func (t *TelegramBot) handleUpdate(update tgbotapi.Update) {
	// Handle messages
	if update.Message != nil {
		// Check if it's a command
		if update.Message.IsCommand() {
			t.handleCommand(update.Message)
		}
	}
}

// handleCommand processes bot commands
func (t *TelegramBot) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		t.handleStartCommand(message)
	case "promote":
		t.handlePromoteCommand(message)
	case "demote":
		t.handleDemoteCommand(message)
	default:
		t.replyToMessage(message, "Unknown command. Available commands: /start, /promote, /demote")
	}
}

// handleStartCommand handles the /start command and registration via deep links
func (t *TelegramBot) handleStartCommand(message *tgbotapi.Message) {
	// Check if this is a deep link with registration hash
	if message.CommandArguments() != "" {
		regHash := message.CommandArguments()
		t.handleRegistration(message, regHash)
	} else {
		t.replyToMessage(message, "Welcome to UTMN Map Bot! I can help you register and authenticate in the UTMN Map system.")
	}
}

// handleRegistration processes registration via deep link
func (t *TelegramBot) handleRegistration(message *tgbotapi.Message, regHash string) {
	ctx := context.Background()

	// Validate registration hash
	err := t.authService.ValidateRegistrationHash(regHash)
	if err != nil {
		var errorMsg string
		switch err {
		case service.ErrHashNotFound:
			errorMsg = "Invalid registration link."
		case service.ErrHashExpired:
			errorMsg = "Registration link has expired. Please request a new one."
		default:
			errorMsg = "An error occurred during registration. Please try again later."
		}
		t.replyToMessage(message, errorMsg)
		return
	}

	// Get user data from message
	telegramID := message.From.ID
	username := message.From.UserName
	if username == "" {
		// If no username, use first name + last name or just ID as string
		if message.From.FirstName != "" {
			username = message.From.FirstName
			if message.From.LastName != "" {
				username += " " + message.From.LastName
			}
		} else {
			username = fmt.Sprintf("user_%d", telegramID)
		}
	} else {
		username = "@" + username
	}

	// Get user photo if available
	var photoURL string
	photos, err := t.bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{
		UserID: message.From.ID,
		Limit:  1,
	})
	if err == nil && photos.TotalCount > 0 {
		photo := photos.Photos[0][0]
		fileInfo, err := t.bot.GetFile(tgbotapi.FileConfig{
			FileID: photo.FileID,
		})
		if err == nil {
			photoURL = fileInfo.Link(t.bot.Token)
		}
	}

	// Register user
	user, err := t.userService.RegisterUser(ctx, telegramID, username, photoURL)
	if err != nil {
		t.replyToMessage(message, "Failed to register. Please try again later.")
		return
	}

	// Consume the hash after successful registration
	if err := t.authService.ConsumeRegistrationHash(regHash); err != nil {
		log.Printf("Failed to consume registration hash: %v", err)
	}

	// Send success message to user
	t.replyToMessage(message, fmt.Sprintf("You have successfully registered as a %s.", user.Role))

	// Send notification to developers chat
	t.notifyDevelopers(fmt.Sprintf("[Registration] %s (ID: %d) has registered.", username, telegramID))
}

// handlePromoteCommand handles the /promote command
func (t *TelegramBot) handlePromoteCommand(message *tgbotapi.Message) {
	ctx := context.Background()

	// Extract target username
	args := message.CommandArguments()
	targetUsername := extractUsername(args)
	if targetUsername == "" {
		t.replyToMessage(message, "Please specify a username to promote. Usage: /promote @username")
		return
	}

	// Get curator from message
	curator, err := t.userService.GetUserByTelegramID(ctx, message.From.ID)
	if err != nil {
		t.replyToMessage(message, "You need to be registered to use this command.")
		return
	}

	// Promote user
	promotedUser, err := t.userService.PromoteUser(ctx, curator.ID, targetUsername)
	if err != nil {
		var errorMsg string
		switch err {
		case userService.ErrUserNotFound:
			errorMsg = fmt.Sprintf("User %s not found.", targetUsername)
		case userService.ErrInsufficientRights:
			errorMsg = "You don't have sufficient rights to promote users."
		default:
			errorMsg = "Failed to promote user. " + err.Error()
		}
		t.replyToMessage(message, errorMsg)
		return
	}

	// Send success message
	t.replyToMessage(message, fmt.Sprintf("%s has been promoted to Admin.", targetUsername))

	// Send notification to developers chat
	t.notifyDevelopers(fmt.Sprintf("[Role Change] Curator %s has promoted %s → Admin", curator.Username, targetUsername))
}

// handleDemoteCommand handles the /demote command
func (t *TelegramBot) handleDemoteCommand(message *tgbotapi.Message) {
	ctx := context.Background()

	// Extract target username
	args := message.CommandArguments()
	targetUsername := extractUsername(args)
	if targetUsername == "" {
		t.replyToMessage(message, "Please specify a username to demote. Usage: /demote @username")
		return
	}

	// Get curator from message
	curator, err := t.userService.GetUserByTelegramID(ctx, message.From.ID)
	if err != nil {
		t.replyToMessage(message, "You need to be registered to use this command.")
		return
	}

	// Demote user
	demotedUser, err := t.userService.DemoteUser(ctx, curator.ID, targetUsername)
	if err != nil {
		var errorMsg string
		switch err {
		case userService.ErrUserNotFound:
			errorMsg = fmt.Sprintf("User %s not found.", targetUsername)
		case userService.ErrInsufficientRights:
			errorMsg = "You don't have sufficient rights to demote users."
		default:
			errorMsg = "Failed to demote user. " + err.Error()
		}
		t.replyToMessage(message, errorMsg)
		return
	}

	// Send success message
	t.replyToMessage(message, fmt.Sprintf("%s has been demoted to User.", targetUsername))

	// Send notification to developers chat
	t.notifyDevelopers(fmt.Sprintf("[Role Change] Curator %s has demoted %s → User", curator.Username, targetUsername))
}

// SendAuthCode sends an authentication code to a user
func (t *TelegramBot) SendAuthCode(ctx context.Context, username string, code string) error {
	// Extract clean username without @
	cleanUsername := strings.TrimPrefix(username, "@")

	// Get user by username
	user, err := t.userService.GetUserByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Send message with code
	msg := tgbotapi.NewMessage(user.TelegramID, fmt.Sprintf("Your entry code is %s (valid for 5 minutes)", code))
	_, err = t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send auth code: %w", err)
	}

	// Send notification to developers chat
	t.notifyDevelopers(fmt.Sprintf("[Auth] %s requested the login code.", username))

	return nil
}

// replyToMessage sends a reply to a message
func (t *TelegramBot) replyToMessage(message *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID
	_, err := t.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send reply: %v", err)
	}
}

// notifyDevelopers sends a message to the developers chat
func (t *TelegramBot) notifyDevelopers(text string) {
	if t.developersChat != 0 {
		msg := tgbotapi.NewMessage(t.developersChat, text)
		_, err := t.bot.Send(msg)
		if err != nil {
			log.Printf("Failed to send notification to developers: %v", err)
		}
	}
}

// extractUsername extracts a username from a string
func extractUsername(text string) string {
	re := regexp.MustCompile(`@(\w+)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return "@" + matches[1]
	}
	return ""
}

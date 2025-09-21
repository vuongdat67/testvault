package security

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "syscall"
    "unicode"
    
    "golang.org/x/term" // replace "golang.org/x/crypto/ssh/terminal"

    "github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/crypto"
)

// PasswordStrength represents password strength level
type PasswordStrength int

const (
    Weak PasswordStrength = iota
    Medium  
    Strong
    VeryStrong
)

func (ps PasswordStrength) String() string {
    switch ps {
    case Weak:
        return "Weak"
    case Medium:
        return "Medium" 
    case Strong:
        return "Strong"
    case VeryStrong:
        return "Very Strong"
    default:
        return "Unknown"
    }
}

// PasswordPolicy defines password requirements
type PasswordPolicy struct {
    MinLength    int
    RequireUpper bool
    RequireLower bool
    RequireDigit bool
    RequireSpecial bool
}

// DefaultPasswordPolicy returns the default password policy
func DefaultPasswordPolicy() PasswordPolicy {
    return PasswordPolicy{
        MinLength:      12,
        RequireUpper:   true,
        RequireLower:   true,
        RequireDigit:   true,
        RequireSpecial: true,
    }
}

// ReadPassword safely reads password from terminal (hidden input)
func ReadPassword(prompt string) (string, error) {
    fmt.Print(prompt)
    
    // Read password without echo
    bytePassword, err := term.ReadPassword(int(syscall.Stdin))
    if err != nil {
        return "", fmt.Errorf("failed to read password: %w", err)
    }
    
    fmt.Println() // New line after hidden input
    password := string(bytePassword)
    
    // Clear sensitive data from memory
    crypto.SecureZero(bytePassword)
    
    return strings.TrimSpace(password), nil
}

// ReadPasswordWithConfirmation reads and confirms password
func ReadPasswordWithConfirmation(prompt string) (string, error) {
    password, err := ReadPassword(prompt + ": ")
    if err != nil {
        return "", err
    }
    
    confirm, err := ReadPassword("Confirm password: ")
    if err != nil {
        return "", err
    }
    
    if password != confirm {
        return "", fmt.Errorf("passwords do not match")
    }
    
    return password, nil
}

// ReadPasswordFromStdin reads password from stdin (for pipes/scripts)
func ReadPasswordFromStdin() (string, error) {
    reader := bufio.NewReader(os.Stdin)
    password, err := reader.ReadString('\n')
    if err != nil {
        return "", fmt.Errorf("failed to read password from stdin: %w", err)
    }
    
    return strings.TrimSpace(password), nil
}

// ValidatePassword checks if password meets policy requirements
func ValidatePassword(password string, policy PasswordPolicy) error {
    if len(password) < policy.MinLength {
        return fmt.Errorf("password must be at least %d characters long", policy.MinLength)
    }
    
    var hasUpper, hasLower, hasDigit, hasSpecial bool
    
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsDigit(char):
            hasDigit = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }
    
    if policy.RequireUpper && !hasUpper {
        return fmt.Errorf("password must contain at least one uppercase letter")
    }
    if policy.RequireLower && !hasLower {
        return fmt.Errorf("password must contain at least one lowercase letter")
    }
    if policy.RequireDigit && !hasDigit {
        return fmt.Errorf("password must contain at least one digit")
    }
    if policy.RequireSpecial && !hasSpecial {
        return fmt.Errorf("password must contain at least one special character")
    }
    
    return nil
}

// CheckPasswordStrength evaluates password strength
func CheckPasswordStrength(password string) PasswordStrength {
    score := 0
    length := len(password)
    
    // Length scoring
    if length >= 8 {
        score++
    }
    if length >= 12 {
        score++
    }
    if length >= 16 {
        score++
    }
    
    // Character variety scoring
    var hasUpper, hasLower, hasDigit, hasSpecial bool
    
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsDigit(char):
            hasDigit = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }
    
    if hasUpper {
        score++
    }
    if hasLower {
        score++
    }
    if hasDigit {
        score++
    }
    if hasSpecial {
        score++
    }
    
    // Convert score to strength
    switch {
    case score >= 7:
        return VeryStrong
    case score >= 5:
        return Strong
    case score >= 3:
        return Medium
    default:
        return Weak
    }
}

// PromptForPasswordWithValidation prompts for password with policy validation
func PromptForPasswordWithValidation(policy PasswordPolicy) (string, error) {
    fmt.Println("Password Requirements:")
    fmt.Printf("- At least %d characters\n", policy.MinLength)
    if policy.RequireUpper {
        fmt.Println("- At least one uppercase letter")
    }
    if policy.RequireLower {
        fmt.Println("- At least one lowercase letter") 
    }
    if policy.RequireDigit {
        fmt.Println("- At least one digit")
    }
    if policy.RequireSpecial {
        fmt.Println("- At least one special character")
    }
    fmt.Println()
    
    for {
        password, err := ReadPasswordWithConfirmation("Enter password")
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }
        
        if err := ValidatePassword(password, policy); err != nil {
            fmt.Printf("Invalid password: %v\n", err)
            continue
        }
        
        strength := CheckPasswordStrength(password)
        fmt.Printf("Password strength: %s\n", strength)
        
        if strength == Weak {
            fmt.Print("Password is weak. Continue anyway? (y/N): ")
            var response string
            fmt.Scanln(&response)
            if strings.ToLower(response) != "y" {
                continue
            }
        }
        
        return password, nil
    }
}
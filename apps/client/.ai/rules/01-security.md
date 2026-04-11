# Security Rules (Client)

## Authentication & Authorization
- Never expose secrets, API keys, tokens, or credentials in code
- Use environment variables for sensitive configuration
- Implement proper authentication on all protected routes
- Use secure storage for tokens (httpOnly cookies preferred)

## Input Validation
- Validate and sanitize all user inputs
- Escape output to prevent XSS attacks
- Use React's built-in XSS protection

## Data Protection
- Never log sensitive information (passwords, tokens, PII)
- Use secure connections (HTTPS)
- Implement proper session management

## Dependencies
- Keep dependencies updated to patch security vulnerabilities
- Audit dependencies for known vulnerabilities regularly

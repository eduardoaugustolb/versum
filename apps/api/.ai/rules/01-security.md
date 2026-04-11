# Security Rules (API)

## Authentication & Authorization
- Never expose secrets, API keys, tokens, or credentials in code
- Use environment variables for sensitive configuration
- Implement proper authentication on all protected endpoints
- Follow principle of least privilege for permissions

## Input Validation
- Validate and sanitize all user inputs
- Use parameterized queries to prevent SQL injection
- Implement rate limiting
- Implement CSRF protection on state-changing operations

## Data Protection
- Never log sensitive information (passwords, tokens, PII)
- Encrypt sensitive data at rest and in transit
- Use secure hashing algorithms for passwords (bcrypt, argon2)
- Implement proper session management

## Dependencies
- Keep dependencies updated to patch security vulnerabilities
- Audit dependencies for known vulnerabilities regularly

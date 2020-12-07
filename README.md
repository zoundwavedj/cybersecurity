# MoneyLion Cybersecurity Assessment Demo

## How to
- Download

## Assumptions

- Network packet transfer security is covered externally (eg. SSL)

## Notable pitfalls
  
- Manual garbage collection of tokens will cause accessTokens to pileup if `/refresh` is called before it's expiry
- Adding on to that, logging out and logging back in repeatedly will also cause a pile up of refresh tokens :)

## Ideas For Improvements

- Add request logging for easier monitoring of endpoints being called
- If user's password is hashed with an older version of Argon2, make prompt to let user update password (ex. password expiry mechanism)
- Add user/account locking mechanism upon multiple login failures
- Use a cache with auto drop/expiry mechanism for token management instead of handling deletion manually (ex. Redis)
- Standardize request/response objects with their domain counterpart (eg. User) for easier development and debugging
- Use a proper state management framework on frontend (or maybe just handle hooks properly heh)
- Add input validations on both backend/frontend (ex. dob format, email format, etc)

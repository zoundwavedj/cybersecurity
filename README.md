# MoneyLion Cybersecurity Assessment Demo

## Prerequisites

- NodeJS v8+
- Go (optional)

## How to run
  
### Frontend

- Clone the repository
- Navigate to the frontend directory
- Run `npm install` or `yarn install` to install dependencies
- Run `npm start` or `yarn start`

### Backend

- Download the binary from the releases page
- Ensure you have these env vars exported
> "SUPERUSERNAME": `<string>`,
> "ACCESS_SECRET": `<string>`,
> "REFRESH_SECRET": `<string>`,
> "ENCRYPT_KEY": `<32bytes string in hex format (64 chars)>`
- If you opt to build/run it yourself
- Clone the repository
- Ensure you have `gcc` installed (via `build-essentials` on mac, `MSYS2` on windows)
- Run `go run .` from the root directory

## Notes

- For the sake of simplicity, all data stores are using a single sqlite file. A proper alternative would be MySQL for persistent storage, and Redis for volatile storage
- Frontend code is also super messy to keep things 'simple' although that's quite counter-intuitive since putting a bunch of components together in a single file introduces more room for error :)

## Assumptions

- Network packet transfer security is covered externally (eg. SSL)

## Notable pitfalls
  
- Manual garbage collection of tokens will cause accessTokens to pileup if `/refresh` is called before it's expiry
- Adding on to that, logging out and logging back in repeatedly will also cause a pile up of refresh tokens :)

## Ideas for improvements

- Add request logging for easier monitoring of endpoints being called
- If user's password is hashed with an older version of Argon2, make prompt to let user update password (ex. password expiry mechanism)
- Add user/account locking mechanism upon multiple login failures
- Use a cache with auto drop/expiry mechanism for token management instead of handling deletion manually (ex. Redis)
- Standardize request/response objects with their domain counterpart (eg. User) for easier development and debugging
- Use a proper state management framework on frontend (or maybe just handle hooks properly heh)
- Add input validations on both backend/frontend (ex. dob format, email format, etc)

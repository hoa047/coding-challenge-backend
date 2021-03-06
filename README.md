Catalog of stolen bikes
====

![](https://images.unsplash.com/photo-1556316384-12c35d30afa4?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=3450&q=80)

Stolen bikes is a typical problem in Malmö, where the Docly HQ is. We need your help to build a service that will help the local police to solve bike thefts in the area.

# Requirements

- [ ] The police wants to able to add, edit and remove officers. (See data model suggestion below)
- [ ] Private citizens want to be able to report stolen bikes. (See data model suggestion below)
- [ ] The system should assign a police offers to handle stolen bike cases as they are being reported.
  - [ ] A police officer can only handle one case at a time.
- [ ] The police should be able to report bikes as found.
  - [ ] When the police finds a bike the case should be marked as resolved and the police officer in charge of the case should be marked as available to take a new case.
- [ ] The system should be able to assign unassigned cases to police officers as they become available.
- [ ] Police and private citizens both want to be able to list the reported bike thefts and their status.

# Tech requirements

- Go or Node.js (Typescript is a plus)
- Tests 
- You may use any database

# Your challenge

- Create an API that satisfies the requirments above.
- Create documentation on how to use Postman (or similar tool) to interact with the API.
- You can use any boilerplate and tools that you want to but we advice you to keep it simple. A clean, robust API is what we're looking for, and we're usually in favor of using all available tools and tricks to get things done.
- We prefer Docker to run our services, but as long as you have clear instructions on how to run yours you may use whatever you want.

# Instructions

- Fork this repo
- Build a clean and robust API
- Let us know that you've completed the challenge and how we can test it.

# Suggested data model

### Police officers

| Column | Type    | Description     |
| ---    | ---     | ---             |
| id     | Integer |                 |
| name   | String  | Name of officer |

### Bike thefts

| Column      | Type     | Description                              |
| ---         | ---      | ---                                      |
| id          | Integer  |                                          |
| title       | String   | Title of report                          |
| brand       | String   | Brand of bike                            |
| city        | String   | City where theft occured                 |
| description | String   | Description of bike                      |
| reported    | DateTime | Date and time when theft was reported    |
| updated     | DateTime | Date and time when case was last updated |
| solved      | Boolean  | True when case has been solved           |
| officer     | Integer  | Officer in charge of the case            |
| image       | Blob     | Image of bike                            |


# License

This project is licensed under MIT. Feel free to use it anyway you see fit.

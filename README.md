# Pento tech challenge

Thanks for taking the time to do our tech challenge.

The challenge is to build a small full stack web app, that can help a freelancer track their time.

It should satisfy these user stories:

- As a user, I want to be able to start a time tracking session
- As a user, I want to be able to stop a time tracking session
- As a user, I want to be able to name my time tracking session
- As a user, I want to be able to save my time tracking session when I am done with it
- As a user, I want an overview of my sessions for the day, week and month
- As a user, I want to be able to close my browser and shut down my computer and still have my sessions visible to me when I power it up again.

## Getting started

You can fork this repo and use the fork as a basis for your project. We don't have any requirements on what stack you use to solve the task, so there is nothing set up beforehand.

## Timing

- Don't spend more than a days work on this challenge. We're not looking for perfection, rather try to show us something special and have reasons for your decisions.
- Get back to us when you have a timeline for when you are done.

## Notes

- Please focus on code quality and showcasing your skills regarding the role you are applying to.

## Development notes

The time tracking session mentioned above is named "TimeTracker" or "Tracker" entity.

The code follows Clean Architecture framework.

The tests aren't fully covered due to time constraint since it would implicate more than a days work, however there are some integration tests for the repository layer.

The core API exposes full CRUD operations for the Tracker entity.

There are some meta information in the domain that adds created_at, updated_at, a soft delete flag and a version(even though that there isn't going to be concurrent applications operating over the entities, just a nice to have).

All the routes are available on:

Fetch a Tracker

GET /api/v1/tracker/{id}

List Trackers

GET /api/v1/tracker

You can further add start_date and end_date as query parameters to the request to get a set of trackers which were created between the timestamps.

GET /api/v1/tracker/?start_date={rfc3339_timestamp}&end_date={rfc3339_timestamp}

Create tracker

POST /api/v1/tracker

Update Tracker

PUT /api/v1/tracker/{id}

Delete Tracker

DELETE /api/v1/tracker/{id}


## Running tests

There are some integration tests that can be run, make sure to run docker-compose up -d before-hand.

go test -tags=integrationdb -v -p=1 ./backend/...

## Starting API

docker-compose up -d

The API is exposed at port 8080.

## Starting frontend app

In frontend folder:

npm start

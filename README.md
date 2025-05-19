# (Senior) Fullstack Engineer (m/w/d) @OofOne Takehome

As part of our application process, we'd like to see you approach to technical challenges by giving you a small
assignment that resembels the challenges awaiting you once you join OofOne. It should take you no more than a few hours
to complete the assignment, but any extra polish or features you might want to put in will not go unnoticed.

## [](https://github.com/OofOne-SE/senior-software-engineer-takehome#the-assignment)The assignment

You will find two files in this repository (besides this `README`): `columns.yaml` and `weather.dat`.
`weather.dat` contains mock weather data collected over several months. `columns.yaml` contains the headers for the
given data.

Your task is to develop a small API using Go that works with the data as if it was realtime data.
Write a small programm in the scripting language of your choice that iterates over the rows in the data file and sends
them to an endpoint of yours one by one.
The endpoint should accept the raw data and store it in a structured manner.

Please also create endpoints to:

- Retrieve the weather data for a given day
- Retrieve the weather data for a range of days
- Expose a websocket connection that transmits the latest data

Write tests as you find necessary and add simple documentation.
Use any database.

Consider the case that the data stream might increase in frequency in the future and the application will need to store
larger amounts of data.

## [](https://github.com/OofOne-SE/senior-software-engineer-takehome#requirements)Requirements

You may choose whatever technologies you prefer, the only requirement is Go as the backend language.

If you have any questions, please ask!

To complete your takehome, please fork this repo and commit your work to your fork. When you are ready for us to look at
it, give us access to your fork so we can review and run it.

---

# Documentation

The application consists of three main components:

1. [API](#api)
2. [Frontend](#frontend)
3. [Seeder](#seeder)

## API

The API is built using Go and exposes the following endpoints:

- `POST /weather`: Accepts raw weather data and stores it in a structured manner.
- `GET /weather/day`: Retrieves weather data for a given date. The date should be provided in the format `YYYY-MM-DD` as
  a query parameter.
- `GET /weather/range`: Retrieves weather data for a range of dates. The start and end dates should be provided as query
  parameters.
- `GET /ws`: Broadcasts the latest weather data over a WebSocket connection.

*Thoughts:*

- When the project becomes bigger, one could add swagger to document the API and use the swagger files to generate the
  client code for the API.
- The Websocket connection is very basic. It would be more scalable to return an object that specifies the type
  (especially when there is not only weather information sent) and the data.
- The API is not secured. It could be secured using e.g. an API key to restrict the access.

## Frontend

The frontend is served by the same go server as the API. It is built on HTML, CSS (TailwindCSS), and JavaScript.
It listens to the WebSocket connection and updates the UI with the latest weather data.

*Thoughts:*

- The frontend could also be built using a scalable framework like Vue.js.
- It could show the weather data via a chart.
- The frontend could be separated from the API and served by a different light weight webserver to make the parts of the
  application more independent and scalable. The frontend could be even served by a CDN.
- The frontend could also serve as a client for the whole api and not only the websocket connection.

## Seeder

The seeder is a small application that reads the `weather.dat` file and sends the data to the API endpoint.
One can configure the seeder to send the data with a specific delay between each request. The default delay is 2 second.

*Thoughts:*

- There is a dependency between the seeder delay and an animation in the Frontend. If the delay is too short, the
  animation duration will be too long and the bars will not display the data correctly.
- The seeder could be configured via an `.env` file (for now the configuration is played in `seeder/index.js`.
- To make the seeder more robust, one could add a retry mechanism in case the API is not reachable.
- The seeder is built using JavaScript, it could be built in TypeScript (or even Go) to make it more robust.


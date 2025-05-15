import fs from 'fs';

const logger = console

const config = {
    lineDelimiter: "\n",
    columnDelimiter: "\t",
    weatherDataFilePath: "../weather.dat",
}


const loadWeatherData = (lineDelimiter, columnDelimiter, filePath) => {
    console.log(filePath);
    const weatherDataFile = fs.readFileSync(filePath, 'utf8')
    return weatherDataFile.split(lineDelimiter).filter(line => !!line).map(line => {
        const [date, humidity, temperature] = line.split(columnDelimiter)
        return {
            date: date, humidity: parseFloat(humidity), temperature: parseFloat(temperature)
        }
    })
}

const isWeatherDataValid = (weatherData) => {
    if (!weatherData.length) {
        logger.warn('Weather data is empty.')
        return false
    }
    const isRecordInvalid = record => {
        const isInvalid = !record.date || isNaN(record.humidity) || isNaN(record.temperature)
        if (isInvalid) {
            logger.warn(`Invalid record: ${JSON.stringify(record)}`)
        }
        return isInvalid
    }

    if (weatherData.some(isRecordInvalid)) {
        logger.warn('Weather data contains invalid records.')
        return false
    }

    return true
}

const getWeatherDataAndSendToApi = async ({ lineDelimiter, columnDelimiter, weatherDataFilePath }) => {
    const weatherData = loadWeatherData(lineDelimiter, columnDelimiter, weatherDataFilePath)
    if (!isWeatherDataValid(weatherData)) {
        logger.warn('Weather data is invalid.')
        return
    }
    // todo: send weather data to API
    console.log("done");
}

getWeatherDataAndSendToApi(config);

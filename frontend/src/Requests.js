const headers = {
    mode: "cors",
    cache: "no-cache",
    headers: {
        "Content-Type": "application/json; charset=utf-8"
    },
    referrerPolicy: "no-referrer"
}

const url = "http://localhost:8080"

const getTrackerByID = async (id) => {
    const uri = `${url}/tracker/${id}`;
    const requestOptions = {...headers, method: "GET"};
    const response = await fetch(uri, requestOptions);

    return await response.json();
}

const getTrackers = async (start, end) => {
    let uri = "";
    start !== undefined && end !== undefined ? uri = `http://localhost:8080/tracker?start_date=${start}&end_date=${end}` : uri = `http://localhost:8080/tracker`;
    const requestOptions = {...headers, method: "GET"};
    
    const response = await fetch(uri, requestOptions);

    return await response.json();
}

const createTracker = async (body) => {
    const uri = `${url}/tracker`;
    const requestOptions = {
        ...headers,
        method: "POST",
        body: JSON.stringify(body)
    };
    const response = await fetch(uri, requestOptions);

    return await response.json();
}

const updateTracker = async (body, id) => {
    const uri = `${url}/tracker/${id}`;
    const requestOptions = {
        ...headers,
        method: "PUT",
        body: JSON.stringify(body)
    };
    const response = await fetch(uri, requestOptions);

    return await response.json();
}

const deleteTracker = async (id) => {
    const uri = `${url}/tracker/${id}`;
    const requestOptions = {...headers, method: "DELETE"};
    const response = await fetch(uri, requestOptions);

    return await response.text();
}

module.exports = {
    getTrackerByID: getTrackerByID,
    getTrackers: getTrackers,
    createTracker: createTracker,
    updateTracker: updateTracker,
    deleteTracker: deleteTracker,
}
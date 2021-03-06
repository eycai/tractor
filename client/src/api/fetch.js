const _apiHost = "http://localhost:3000/api";

async function request(url, params, method = "GET") {
  const options = {
    method,
    headers: {
      "Content-Type": "application/json",
    },
  };

  // if params exists and method is GET, add query string to url
  // otherwise, just add params as a "body" property to the options object
  if (params) {
    if (method === "GET") {
      url += "?" + objectToQueryString(params);
    } else {
      options.body = JSON.stringify(params); // body should match Content-Type in headers option
    }
  }

  const response = await fetch(_apiHost + url, options);

  // show an error if the status code is not 200
  if (response.status !== 200) {
    return generateErrorResponse(
      response.status,
      "The server responded with an unexpected status."
    );
  }

  let result = null;
  console.log(response);
  try {
    result = await response.json();
  } catch (e) {
    console.log(e);
    return generateErrorResponse(response.status, "Couldn't jsonify");
  }
  return { status: 200, payload: result };
}

function generateErrorResponse(status, message) {
  return {
    status: status,
    payload: message,
  };
}

function objectToQueryString(obj) {
  return Object.keys(obj)
    .map((key) => key + "=" + obj[key])
    .join("&");
}

function get(url, params) {
  return request(url, params);
}

function post(url, params) {
  return request(url, params, "POST");
}

export { get, post };

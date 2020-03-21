const _apiHost = 'http://localhost:8080';

async function request(url, params, method = 'GET') {
  const options = {
    method,
    headers: {
      'Content-Type': 'application/json'
    }
  };

  // if params exists and method is GET, add query string to url
  // otherwise, just add params as a "body" property to the options object
  if (params) {
    if (method === 'GET') {
      url += '?' + objectToQueryString(params);
    } else {
      options.body = JSON.stringify(params); // body should match Content-Type in headers option
    }
  }

  const response = await fetch(_apiHost + url, options);

  // show an error if the status code is not 200
  if (response.status !== 200) {
    return generateErrorResponse(
      'The server responded with an unexpected status.'
    );
  }

  const result = await response.json();
  return result;
}

function generateErrorResponse(message) {
  return {
    status: 'error',
    message
  };
}

function objectToQueryString(obj) {
  return Object.keys(obj)
    .map(key => key + '=' + obj[key])
    .join('&');
}

function get(url, params) {
  return request(url, params);
}

function create(url, params) {
  return request(url, params, 'POST');
}

function update(url, params) {
  return request(url, params, 'PUT');
}

function remove(url, params) {
  return request(url, params, 'DELETE');
}

export default {
  get,
  create,
  update,
  remove
};
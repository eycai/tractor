import Fetch from './fetch.js';

// GET
async function getHelloWorld() {
  const resp = await Fetch.get('/hello');
  return resp;
}

export default {
  getHelloWorld
};

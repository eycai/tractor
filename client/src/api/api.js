import Fetch from './fetch.js';

// GET
async function getHelloWorld() {
  const resp = await Fetch.get('/hello');
  return resp;
}

export default {
  getHelloWorld
};

// POST
// async function createBook() {
//   const request = await Fetch.create('/books', {
//     title: 'Code and Other Laws of Cyberspace',
//     author: 'Lawrence Lessig'
//   });
// }

// // PUT
// async function updateBook(bookId) {
//   const request = await Fetch.update('/books/' + bookId, {
//     title: 'How to Live on Mars',
//     author: 'Elon Musk'
//   });
// }

// // DELETE
// async function removeBook(bookId) {
//   const request = await Fetch.remove('/books/' + bookId);
// }

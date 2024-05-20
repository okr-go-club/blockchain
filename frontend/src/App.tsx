import React from 'react';
import './App.css';

function MyButton() {
  return (
    <button>I'm a button</button>
  );
}

export default function App() {
  return (
    <div>
      <h1>Welcome to my app</h1>
      <MyButton />
    </div>
  );
}

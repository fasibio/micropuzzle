import React from 'react'
import ReactDOM from 'react-dom'
import App from './App'
import { registerCustomElement } from './customElementRegister'

const Root = () => {
  return (
  <React.StrictMode>
        <App />
  </React.StrictMode>)
 }

registerCustomElement({
  attributes: [],
  component: Root,
  name: 'header-component',
})
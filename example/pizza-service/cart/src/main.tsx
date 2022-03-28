import React, { useEffect, useState } from 'react'
import ReactDOM from 'react-dom'
import App from './App'
import { registerCustomElement } from './customElementRegister'

const Root = () => {
  const [t, setT] = useState(0) 
  useEffect(() => {
    setT(t + 1)
  }, []);
  return (
  <React.StrictMode>
        {t}
        <App />
  </React.StrictMode>)
}
registerCustomElement({
  attributes: [],
  component: Root,
  name: 'cart-component',
})
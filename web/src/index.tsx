import * as React from 'react'
import * as ReactDOM from 'react-dom'
import { HashRouter as Router } from 'react-router-dom'

import App from './App'

import './style/index.less'

ReactDOM.render(
  <Router><App/></Router>,
  document.getElementById('root') as HTMLElement,
)

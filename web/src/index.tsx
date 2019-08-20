import * as React from 'react'
import * as ReactDOM from 'react-dom'
import { history } from './utils/URLUtil'
import { Router } from 'react-router-dom'

import App from './App'

import './style/index.less'

ReactDOM.render(
  <Router history={history}><App/></Router>,
  document.getElementById('root') as HTMLElement,
)

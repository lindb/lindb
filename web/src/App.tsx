import Content from 'containers/layout/Content'
import Login from 'containers/Login/Login'
import * as React from 'react'
import { Route, Switch } from 'react-router-dom'

interface AppProps {
}

interface AppState {
}

export default class App extends React.Component<AppProps, AppState> {

  constructor(props: Readonly<AppProps>) {
    super(props)
    this.state = { collapsed: false }
  }

  public render() {
    return (
      <Switch>
        <Route exact={true} path="/login" component={Login}/>
        <Route exact={false} path="/" component={Content}/>
      </Switch>
    )
  }
}
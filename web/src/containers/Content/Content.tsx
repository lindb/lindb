import { Layout } from 'antd'
import * as React from 'react'
import { Redirect, Route, Switch } from 'react-router-dom'

import Header from '../../components/Layout/Header'
import Footer from '../../components/Layout/Footer'
import SiderMenu from '../../components/Layout/SiderMenu'
import ChartTooltip from '../../components/Chart/ChartTooltip'

import SearchPage from '../Search'
import HomePage from '../Home/Home'
import MonitoringSystem from '../Monitoring/System'

const { Content: AntDContent } = Layout

interface ContentProps {
}

interface ContentStatus {
}

export default class Content extends React.Component<ContentProps, ContentStatus> {
  constructor(props: ContentProps) {
    super(props)
    this.state = {}
  }

  render() {
    return (
      <Layout className="lindb-sider-layout">
        {/* Sider Bar Menu */}
        <SiderMenu/>

        {/* Content Area */}
        <Layout>
          <AntDContent className="lindb-content-container">
            <Header/>

            <ChartTooltip/>

            <Switch>
              <Route exact={true} path="/" component={HomePage}/>
              <Route exact={true} path="/search" component={SearchPage}/>
              <Route exact={true} path="/monitoring/system" component={MonitoringSystem}/>
            </Switch>
          </AntDContent>

          <Footer sider={true}/>
        </Layout>
      </Layout>
    )
  }
}
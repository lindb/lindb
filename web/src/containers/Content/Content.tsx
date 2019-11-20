import { Layout } from 'antd'
import * as React from 'react'
import { Route, Switch } from 'react-router-dom'

import Header from '../../components/Layout/Header'
import Footer from '../../components/Layout/Footer'
import SiderMenu from '../../components/Layout/SiderMenu'
import ChartTooltip from '../../components/Chart/ChartTooltip'

import SearchPage from '../Search/DataSearch'
import OverviewPage from '../Overview/Overview'
import StorageClusterDetailPage from '../Overview/StorageClusterDetail'
import Cluster from '../Setting/Cluster'
import Database from '../Setting/Database'
import MonitoringSystem from '../Monitor/System'
import { CreateDatabaseWithForm } from '../Setting/NewDatabase'

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
        <Layout className="lindb-layout">
          <AntDContent className="lindb-content-container">
            <Header/>

            <ChartTooltip/>

            <Switch>
              <Route exact={true} path="/" component={OverviewPage}/>
              <Route exact={true} path="/storage/cluster/:clusterName" component={StorageClusterDetailPage}/>
              <Route exact={true} path="/search" component={SearchPage}/>
              <Route exact={true} path="/monitoring/system" component={MonitoringSystem}/>
              <Route exact={true} path="/setting/cluster" component={Cluster}/>
              <Route exact={true} path="/setting/database" component={Database}/>
              <Route exact={true} path="/setting/database/new" component={CreateDatabaseWithForm}/>
            </Switch>
          </AntDContent>

          <Footer sider={true}/>
        </Layout>
      </Layout>
    )
  }
}
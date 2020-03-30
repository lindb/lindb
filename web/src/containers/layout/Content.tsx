import { Layout } from 'antd'
import * as React from 'react'
import { Route, Switch } from 'react-router-dom'
import ChartTooltip from '../../components/metric/ChartTooltip'
import Footer from './Footer'
import Header from './Header'
import SiderMenu from './SiderMenu'
import Database from '../admin/Database'
import Storage from '../admin/Storage'
import Runtime from '../monitoring/Runtime'
import OverviewPage from '../home/Overview'
import StorageClusterDetailPage from '../home/StorageClusterDetail'
import SearchPage from '../query/MetricDataSearch'

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
              <Route exact={true} path="/monitoring/runtime" component={Runtime}/>
              <Route exact={true} path="/admin/storage" component={Storage}/>
              <Route exact={true} path="/admin/database" component={Database}/>
            </Switch>
          </AntDContent>

          <Footer sider={true}/>
        </Layout>
      </Layout>
    )
  }
}
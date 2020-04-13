import { HomeOutlined } from '@ant-design/icons';
import { Breadcrumb } from 'antd';
import { observer } from 'mobx-react';
import { BreadcrumbStatus } from 'model/Breadcrumb';
import * as React from 'react';
import { Link, withRouter } from 'react-router-dom';
import StoreManager from 'store/StoreManager';

interface BreadcrumbHeaderProps {
  location: any
}

interface BreadcrumbHeaderStatus {
}

@observer
class BreadcrumbHeader extends React.Component<BreadcrumbHeaderProps, BreadcrumbHeaderStatus> {
  breadcrumbStore: any

  constructor(props: BreadcrumbHeaderProps) {
    super(props)

    this.breadcrumbStore = StoreManager.BreadcrumbStore
  }

  render() {
    const breadcrumbItems = this.breadcrumbStore.breadcrumbList.map((item: BreadcrumbStatus) => {
      return <Breadcrumb.Item key={item.label}><Link to={item.path}>{item.label}</Link></Breadcrumb.Item>
    })
    return (
      <div className="lindb-header__breadcrumb">
        {this.breadcrumbStore.breadcrumbList.length > 0 && (<HomeOutlined />)}
        <Breadcrumb>
          {breadcrumbItems}
        </Breadcrumb>
      </div>
    )
  }
}

// @ts-ignore
export default withRouter(BreadcrumbHeader)
import * as React from 'react'
import { MENUS } from '../config/menu'
import { Breadcrumb, Icon } from 'antd'
import { withRouter, Link } from 'react-router-dom'

import defaults from 'lodash-es/defaults'
import flattenDeep from 'lodash-es/flattenDeep'

interface BreadcrumbHeaderProps {
  location: any
}

interface BreadcrumbHeaderStatus {
}

class BreadcrumbHeader extends React.Component<BreadcrumbHeaderProps, BreadcrumbHeaderStatus> {
  breadcrumbNameMap: object

  constructor(props: BreadcrumbHeaderProps) {
    super(props)
    this.state = {}

    this.breadcrumbNameMap = defaults(...flattenDeep(MENUS.map(item => {
      if (item.children) {
        return [
          { [ item.path ]: item.title },
          ...item.children.map(child => ({ [ item.path + child.path ]: child.title })),
        ]
      }
      return { [ item.path ]: item.title }
    })))
  }

  render() {
    const { location } = this.props
    const pathSnippets = location.pathname === '/' ? [ '' ] : location.pathname.split('/').filter(i => i)
    const breadcrumbItems = pathSnippets.map((_, index) => {
      const url = `/${pathSnippets.slice(0, index + 1).join('/')}`
      const title = this.breadcrumbNameMap[ url ]

      return title ? (<Breadcrumb.Item key={url}>{title}</Breadcrumb.Item>) : null
    }).filter(Boolean)

    return (
      <div className="lindb-header__breadcrumb">
        {breadcrumbItems.length > 0 && (<Icon type="compass" />)}
        <Breadcrumb>
          {breadcrumbItems}
        </Breadcrumb>
      </div>
    )
  }
}

export default withRouter(BreadcrumbHeader)
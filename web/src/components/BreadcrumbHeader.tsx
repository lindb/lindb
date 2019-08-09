import * as React from 'react'
import { MENUS } from '../config/menu'
import { Breadcrumb, Icon } from 'antd'
import { withRouter } from 'react-router-dom'

interface BreadcrumbHeaderProps {
  location: any
}

interface BreadcrumbHeaderStatus {
}

class BreadcrumbHeader extends React.Component<BreadcrumbHeaderProps, BreadcrumbHeaderStatus> {
  breadcrumbNameMap: any

  constructor(props: BreadcrumbHeaderProps) {
    super(props)

    const menu = MENUS.map(item => {
      if (item.children) {
        return [
          { [ item.path ]: item.title },
          ...item.children.map(child => ({
            [ item.path + child.path ]: child.title
          })),
        ]
      }
      return [{ [ item.path ]: item.title }]
    })

    this.breadcrumbNameMap = {}
    menu.forEach(m => {
      Object.assign(this.breadcrumbNameMap, ...m)
    })
  }

  render() {
    const { location } = this.props
    const pathSnippets = location.pathname === '/' ? [ '' ] : location.pathname.split('/').filter(Boolean)
    const breadcrumbItems = pathSnippets.map((_: any, index: number) => {
      const url = `/${pathSnippets.slice(0, index + 1).join('/')}`
      const title = this.breadcrumbNameMap[url]
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
// @ts-ignore
export default withRouter(BreadcrumbHeader)
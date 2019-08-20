import * as React from 'react'
import { Link } from 'react-router-dom'
import { Layout, Icon, Menu } from 'antd'

import { MENUS } from '../../config/menu'

const { Sider } = Layout
const Logo = require('../../assets/images/logo_title.png')

interface SiderMenuProps {
}

interface SiderMenuStatus {
}

export default class SiderMenu extends React.Component<SiderMenuProps, SiderMenuStatus> {
  renderMenu(menus: Array<any>, parentPath?: string) {
    const IconTitle = (icon: string, title: string) => <span><Icon type={icon}/>{title}</span>

    return menus.map(menu => {
      const path = parentPath ? parentPath + menu.path : menu.path

      return menu.children
        ? (
          <Menu.SubMenu key={menu.path} title={IconTitle(menu.icon, menu.title)}>
            {this.renderMenu(menu.children, menu.path)}
          </Menu.SubMenu>
        )
        : (
          <Menu.Item key={path}>
            <Link to={path}>{IconTitle(menu.icon, menu.title)}</Link>
          </Menu.Item>
        )
    })
  }

  render() {
    const { location: { hash } } = window
    const path = hash.replace('#', '')

    return (
      <Sider className="lindb-sider" collapsible={true} trigger={null}>
        {/* Logo */}
        <div className="lindb-sider__logo">
          <Link to="/"><img src={Logo} alt="LinDB"/></Link>
        </div>

        {/* Menu */}
        <Menu
          mode="inline"
          className="lindb-sider__menu"
          defaultOpenKeys={[ '/monitoring', '/setting' ]}
          selectedKeys={[ path ]}
        >
          {this.renderMenu(MENUS)}
        </Menu>
      </Sider>
    )
  }
}
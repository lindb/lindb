import * as React from 'react'
import { Icon, Layout, Menu, Dropdown, Divider } from 'antd'
import { Link } from 'react-router-dom'
import BreadcrumbHeader from '../BreadcrumbHeader'
import { LOCALSTORAGE_TOKEN } from '../../config/config'
import { redirectTo } from '../../utils/URLUtil'

const { Header: AntDHeader } = Layout

interface HeaderProps {
}

interface HeaderStatus {
}

export default class Header extends React.Component<HeaderProps, HeaderStatus> {
  handleLogout() {
    localStorage.removeItem(LOCALSTORAGE_TOKEN)
    redirectTo(window.location.pathname)
  }

  render() {
    const isLogin = !!localStorage.getItem(LOCALSTORAGE_TOKEN)

    const logout = (
      <Menu><Menu.Item><span onClick={this.handleLogout}><Icon type="logout" /> Logout</span></Menu.Item></Menu>
    )

    const user = isLogin
      ? <Dropdown overlay={logout} placement="bottomCenter"><span><Icon type="user" />Admin</span></Dropdown>
      : <Link to="/login">Login</Link>

    return (
      <AntDHeader className="lindb-header">
        <BreadcrumbHeader/>

        <ul className="lindb-header__menu">
          <li className="lindb-header__menu-item bold">{user}</li>
          <div style={{float: 'right'}}><Divider type="vertical" /></div>
          <li className="lindb-header__menu-item">
            <a href="https://lindb.io" rel="noopener noreferrer" target="_blank"><Icon type="question-circle" />Help</a>
          </li>
          <li className="lindb-header__menu-item">
            <a href="https://github.com/lindb/lindb" rel="noopener noreferrer" target="_blank"><Icon type="github" />GitHub</a>
          </li>
        </ul>
      </AntDHeader>
    )
  }
}
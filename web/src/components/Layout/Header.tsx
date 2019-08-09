import * as React from 'react'
import { Icon, Layout } from 'antd'
import { Link } from 'react-router-dom'
import BreadcrumbHeader from '../BreadcrumbHeader'

const { Header: AntDHeader } = Layout

interface HeaderProps {
}

interface HeaderStatus {
}

export default class Header extends React.Component<HeaderProps, HeaderStatus> {
  render() {
    return (
      <AntDHeader className="lindb-header">
        <BreadcrumbHeader/>

        <ul className="lindb-header__menu">
          <li className="lindb-header__menu-item bold">
            <Link to="/login">Login</Link>
          </li>
          <li className="lindb-header__menu-item">
            <a href="https://github.com/eleme/lindb/wiki" rel="noopener noreferrer" target="_blank"><Icon type="read" />Help</a>
          </li>
          <li className="lindb-header__menu-item">
            <a href="https://github.com/eleme/lindb" rel="noopener noreferrer" target="_blank"><Icon type="github" />GitHub</a>
          </li>
        </ul>
      </AntDHeader>
    )
  }
}
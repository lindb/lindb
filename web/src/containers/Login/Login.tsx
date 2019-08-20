import * as React from 'react'
import { Link } from 'react-router-dom'
import { autobind } from 'core-decorators'
import { login } from '../../service/login'
import LoginBackground from './LoginBackground'
import Footer from '../../components/Layout/Footer'
import { LOCALSTORAGE_TOKEN } from '../../config/config'
import { Button, Icon, Input, Layout, message } from 'antd'
import { getQueryValueOf, redirectTo } from '../../utils/URLUtil'

const { Content } = Layout
const Logo = require('../../assets/images/logo_title_subtitle.png')

interface LoginProps {
}

interface LoginStatus {
  username: string
  password: string
  isLoading: boolean
}

export default class Login extends React.Component<LoginProps, LoginStatus> {
  constructor(props: LoginProps) {
    super(props)
    this.state = {
      username: '',
      password: '',
      isLoading: false
    }
  }

  @autobind
  handleUsernameChange(e: React.ChangeEvent<HTMLInputElement>) {
    const username = e.target.value
    this.setState({ username })
  }

  @autobind
  handlePasswordChange(e: React.ChangeEvent<HTMLInputElement>) {
    const password = e.target.value
    this.setState({ password })
  }

  @autobind
  async handleLogin() {
    const { username, password } = this.state
    this.setState({ isLoading: true })
    const token = await login(username, password)
    this.setState({ isLoading: false })
    if (token) {
      this.saveToken(token)
    } else {
      message.warning('username or password is incorrect.')
    }
  }

  saveToken(token: string) {
    localStorage.setItem(LOCALSTORAGE_TOKEN, token)

    // redirect
    const from = getQueryValueOf('from')
    setTimeout(() => {redirectTo(from || '/')}, 1000)
  }

  render() {
    const { username, password, isLoading } = this.state

    return (
      <Layout className="lindb-login">
        <LoginBackground/>

        <Content>

          <div className="lindb-login__content">
            <Link to="/"><img className="lindb-login__content__logo" src={Logo} alt=""/></Link>
            <Input
              className="lindb-login__content__input"
              prefix={<Icon type="user" style={{ color: 'rgba(0,0,0,.25)' }}/>}
              placeholder="Enter Username"
              value={username}
              onChange={this.handleUsernameChange}
            />

            <Input.Password
              className="lindb-login__content__input"
              prefix={<Icon type="lock" style={{ color: 'rgba(0,0,0,.25)' }}/>}
              placeholder="Enter Password"
              value={password}
              onChange={this.handlePasswordChange}
            />

            <Button
              className="lindb-login__content__btn"
              type="primary"
              shape="round"
              block={true}
              disabled={isLoading}
              onClick={this.handleLogin}
            >
              {isLoading ? 'Login...' : 'Login'}
            </Button>
          </div>
        </Content>

        <Footer/>
      </Layout>
    )
  }
}
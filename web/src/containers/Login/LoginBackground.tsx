import * as React from 'react'

interface LoginBackgroundProps {
}

interface LoginBackgroundStatus {
}

export default class LoginBackground extends React.Component<LoginBackgroundProps, LoginBackgroundStatus> {
  canvas: React.RefObject<HTMLCanvasElement>

  constructor(props: LoginBackgroundProps) {
    super(props)
    this.state = {}

    this.canvas = React.createRef()
  }

  componentDidMount(): void {
    this.init()
  }

  init() {
    // ...
  }

  render() {
    return (
      <div>
        <canvas ref={this.canvas}/>
      </div>
    )
  }
}
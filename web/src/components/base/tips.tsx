import { Icon } from 'antd'
import * as React from 'react'
import classNames from 'classnames'
import { PREFIXCLS } from '../../config/config'

interface TipsProps {
  tip: string
  icon?: string
  className?: string
  size?: 'large' | 'normal'
}

interface TipsStatus {
}

export default class Tips extends React.Component<TipsProps, TipsStatus> {
  render() {
    const { icon, tip, size, className } = this.props
    const prefixCls = `${PREFIXCLS}-tips`

    const classes = classNames(prefixCls, className, {
      [`${prefixCls}-large`]: size === 'large'
    })

    return (
      <div className={classes}>
        {icon && <span className={`${prefixCls}__icon`}><Icon type="warning"/></span>}
        <span className={`${prefixCls}__content`}>{tip}</span>
      </div>
    )
  }
}
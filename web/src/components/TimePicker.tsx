import * as React from 'react'
import { autobind } from 'core-decorators'
import { TIME_FORMAT } from '../config/config'
import { Button, Col, DatePicker, Divider, Popover, Row } from 'antd'

const RangePicker = DatePicker.RangePicker

interface TimePickerProps {
}

interface TimePickerStatus {
}

export default class TimePicker extends React.Component<TimePickerProps, TimePickerStatus> {
  private timepickerCls = 'lindb-timepicker'

  constructor(props: TimePickerProps) {
    super(props)
    this.state = {}
  }

  /**
   * Quick select button handle
   */
  @autobind
  handleQuickSelectItemClick(value: string) {
    // ...
  }

  /**
   * default reference button
   */
  @autobind
  renderDefaultBtn() {
    return (
      <Button icon="clock-circle">Last 7 Days</Button>
    )
  }

  /**
   * Time picker panel
   * @return {any}
   */
  @autobind
  renderTimePicker() {
    interface QuickSelectItem {
      title: string
      value: string
    }

    const Hours: QuickSelectItem[] = [
      { title: 'Last 30 min', value: 'time > now()-30m' },
      { title: 'Last 1 hour', value: 'time > now()-1h' },
      { title: 'Last 3 hours', value: 'time > now()-3h' },
      { title: 'Last 6 hours', value: 'time > now()-6h' },
      { title: 'Last 12 hours', value: 'time > now()-12h' },
    ]

    const Days: QuickSelectItem[] = [
      { title: 'Last 1 day', value: 'time > now()-1d' },
      { title: 'Last 2 days', value: 'time > now()-2d' },
      { title: 'Last 3 days', value: 'time > now()-3d' },
      { title: 'Last 7 days', value: 'time > now()-7d' },
      { title: 'Last 15 days', value: 'time > now()-15d' },
      { title: 'Last 30 days', value: 'time > now()-30d' },
    ]

    const renderQuickSelectItem = (items: QuickSelectItem[], span: number = 12) => {
      const SelectItems = items.map(item => (
        <Button key={item.value} type="link" block={true} onClick={() => this.handleQuickSelectItemClick(item.value)}>
          {item.title}
        </Button>
      ))

      return <Col span={span}>{SelectItems}</Col>
    }

    const cls = this.timepickerCls

    return (
      <div className={cls}>
        {/* Quick Select */}
        <Divider>Quick Select</Divider>
        <Row>
          {renderQuickSelectItem(Hours)}
          {renderQuickSelectItem(Days)}
        </Row>

        {/* Time Range */}
        <Divider>Time Range</Divider>
        <Row>
          <RangePicker
            className={`${cls}-range-picker`}
            format={TIME_FORMAT}
            showTime={true}
          />
        </Row>
        <Row justify="center">
          <Button type="primary">Apply Range</Button>
        </Row>
      </div>
    )
  }

  render() {
    const timepicker = this.renderTimePicker()
    const defaultBtn = this.renderDefaultBtn()

    return (
      <Popover
        trigger="click"
        placement="bottomRight"
        content={timepicker}
        overlayClassName={`${this.timepickerCls}-popover`}
      >
        {defaultBtn}
      </Popover>
    )
  }
}
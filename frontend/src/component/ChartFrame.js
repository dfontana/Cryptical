import React from 'react';
import { VictoryPie } from 'victory';
import './ChartFrame.css'

class ChartFrame extends React.Component {
  render(){
    return (
      <div id="component-chart">
        <VictoryPie />
      </div>
    )
  }
}

export default ChartFrame;
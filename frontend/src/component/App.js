import React from 'react';
import Sidebar from './Sidebar';
import ChartFrame from './ChartFrame';
import BottomFrame from './BottomFrame';
import Header from './Header'
import './App.css'

class App extends React.Component {
  render(){
    return (
      <div id="component-app">
        <div id="component-app-top">
          <Header />
        </div>
        <div id="component-app-middle">
          <Sidebar />
          <div id="chart"><ChartFrame /></div>
        </div>
        <div id="component-app-bottom">
          <BottomFrame />
        </div>
      </div>
    )
  }
}

export default App;
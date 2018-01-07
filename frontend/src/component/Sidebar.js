import React from 'react';
import './Sidebar.css'

class Sidebar extends React.Component {
  render(){
    return (
      <div id="component-sidebar">
        <button>Bitcoin</button>
        <button>Ethereum</button>
        <button>Litecoin</button>
        <button>IOTA</button>
      </div>
    )
  }
}

export default Sidebar;
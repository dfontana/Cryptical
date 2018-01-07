import React from 'react';
import "./Header.css"

class Header extends React.Component {
  render(){
    return (
      <div id="component-header">
        <div id="component-header-icon">(ICON)</div>
        <div id="component-header-title">Cryptical</div>
      </div>
    )
  }
}

export default Header;
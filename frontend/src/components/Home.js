import React, { Component } from 'react';
import { Spinner } from "@blueprintjs/core";
import "@blueprintjs/core/lib/css/blueprint.css"
import "../styles/Home.css"
import CalendarHeatmap from "react-calendar-heatmap"
import "react-calendar-heatmap/dist/styles.css";

export class Home extends Component {
    constructor(props) {
        super(props);
        this.state = {
            loadingMessage: "",
            data: null
        }
    }

    componentDidMount() {
        window.backend.checkForPathList().then(res => {
            if (res) {
                // no need to scan for repos; can go ahead and generate stats
                this.setState({ loadingMessage: "Generating stats..." });
                window.backend.stats().then(data => {
                    this.setState({ data })
                })
            } else {
                // need to scan for repos
                this.setState({ loadingMessage: "Scanning for git repos..." });
                window.backend.scan().then(() => {
                    this.setState({ loadingMessage: "Generating stats..." });
                    window.backend.stats().then(data => {
                        this.setState({ data })
                    })
                })
            }
        })
    }

    loading = () => {
       return (
           <div className="loading-view">
               <Spinner size={70}/>
               <p className="loading-text">{this.state.loadingMessage}</p>
           </div>
       )
    };

    graph = () => {
        return (
            <div style={{width: 700, marginBottom: 50}}>
                <p className="graph-title">19 contributions in the last six months</p>
                <CalendarHeatmap
                    startDate={new Date('2019-01-01')}
                    endDate={new Date('2019-07-01')}
                    values={[]}
                />
            </div>
        )
    };

    render() {
        return (
            <div className="container">
                {this.state.data === null ? this.loading() : this.graph()}
            </div>
        )
    }
}
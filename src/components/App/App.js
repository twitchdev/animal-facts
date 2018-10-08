import React from 'react'
import Authentication from '../Authentication/Authentication'

import './App.css'

export default class App extends React.Component{
    constructor(props){
        super(props)
        this.Authentication = new Authentication()

        //if the extension is running on twitch or dev rig, set the shorthand here. otherwise, set to null. 
        this.twitch = window.Twitch ? window.Twitch.ext : null
        this.state={
            finishedLoading:false,
            theme:'light',
            animal:'',
            fact:''
        }
    }

    contextUpdate(context, delta){
        if(delta.includes('theme')){
            this.setState(()=>{
                return {theme:context.theme}
            })
        }
    }

    componentDidMount(){
        if(this.twitch){
            this.twitch.onAuthorized((auth)=>{
                this.Authentication.setToken(auth.token, auth.userId)

                if(!this.state.finishedLoading){
                    this.setState(()=>{
                        return {finishedLoading:true}
                    })
                }
            })

            this.twitch.onContext((context,delta)=>{
                this.contextUpdate(context,delta)
            })

            this.twitch.configuration.onChanged(()=>{
                let animal = this.twitch.configuration.broadcaster ? this.twitch.configuration.broadcaster.content : ''
                let fact = this.twitch.configuration.developer ? this.twitch.configuration.developer.content : ''
                
                this.setState(()=>{
                    return{
                        animal,
                        fact
                    }
                })
            })
        }
    }

    componentWillUnmount(){
        if(this.twitch){
            this.twitch.unlisten('broadcast', ()=>console.log('successfully unlistened'))
        }
    }
    
    render(){
        if(this.state.finishedLoading && this.state.animal && this.state.fact){
            return (
                <div className={this.state.theme === 'light' ? "App App-light" : "App App-dark"}>
                    <p>Random fact about {this.state.animal}s!</p>
                    <p>Did you know that: {this.state.fact}</p>
                </div>
            )
        }
        else if(this.state.finishedLoading){
            return(
                <div className={this.state.theme === 'light' ? "App App-light" : "App App-dark"}>
                    Extension not configured.
                </div>
            )
        }
        else{
            return (
                <div className="App">
                    Loading...
                </div>
            )
        }
    }
}
export namespace game {
	
	export class NPCInfo {
	    id: string;
	    name: string;
	    race: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new NPCInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.race = source["race"];
	        this.description = source["description"];
	    }
	}
	export class LocationInfo {
	    id: string;
	    name: string;
	    description: string;
	    npcs: NPCInfo[];
	
	    static createFrom(source: any = {}) {
	        return new LocationInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.npcs = this.convertValues(source["npcs"], NPCInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}


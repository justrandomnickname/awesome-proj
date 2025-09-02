export namespace entities {
	
	export class Point {
	    id: string;
	    name: string;
	    description: string;
	    sub_cluster_id: string;
	    type: string;
	    connections: string[];
	    npcs: string[];
	    is_entry_point: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Point(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.sub_cluster_id = source["sub_cluster_id"];
	        this.type = source["type"];
	        this.connections = source["connections"];
	        this.npcs = source["npcs"];
	        this.is_entry_point = source["is_entry_point"];
	    }
	}
	export class SubCluster {
	    id: string;
	    name: string;
	    description: string;
	    cluster_id: string;
	    entry_points: string[];
	    points: Record<string, Point>;
	
	    static createFrom(source: any = {}) {
	        return new SubCluster(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.cluster_id = source["cluster_id"];
	        this.entry_points = source["entry_points"];
	        this.points = this.convertValues(source["points"], Point, true);
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
	export class Cluster {
	    id: string;
	    name: string;
	    description: string;
	    type: string;
	    main_point: string;
	    sub_clusters: Record<string, SubCluster>;
	    child_clusters?: Record<string, Cluster>;
	
	    static createFrom(source: any = {}) {
	        return new Cluster(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.type = source["type"];
	        this.main_point = source["main_point"];
	        this.sub_clusters = this.convertValues(source["sub_clusters"], SubCluster, true);
	        this.child_clusters = this.convertValues(source["child_clusters"], Cluster, true);
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
	export class Interaction {
	    id: string;
	    type: string;
	    content: string;
	    location_id: string;
	    // Go type: time
	    timestamp: any;
	
	    static createFrom(source: any = {}) {
	        return new Interaction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.content = source["content"];
	        this.location_id = source["location_id"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
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
	export class NPC {
	    id: string;
	    name: string;
	    race: string;
	    location_id: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new NPC(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.race = source["race"];
	        this.location_id = source["location_id"];
	        this.description = source["description"];
	    }
	}
	export class Location {
	    id: string;
	    name: string;
	    description: string;
	    current_state: string;
	    type: string;
	    exits: Record<string, string>;
	    npcs: string[];
	    npcs_detailed?: NPC[];
	    interactions?: Interaction[];
	
	    static createFrom(source: any = {}) {
	        return new Location(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.current_state = source["current_state"];
	        this.type = source["type"];
	        this.exits = source["exits"];
	        this.npcs = source["npcs"];
	        this.npcs_detailed = this.convertValues(source["npcs_detailed"], NPC);
	        this.interactions = this.convertValues(source["interactions"], Interaction);
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
	export class LocationHierarchy {
	    clusters: Record<string, Cluster>;
	
	    static createFrom(source: any = {}) {
	        return new LocationHierarchy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.clusters = this.convertValues(source["clusters"], Cluster, true);
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
	
	
	export class SaveInfo {
	    name: string;
	    // Go type: time
	    created_at: any;
	    filename: string;
	
	    static createFrom(source: any = {}) {
	        return new SaveInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.filename = source["filename"];
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


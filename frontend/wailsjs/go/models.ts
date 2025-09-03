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
	export class ConnectionInfo {
	    id: string;
	    name: string;
	    description: string;
	    sub_cluster_id: string;
	    type: string;
	    connections: string[];
	    connection_names: Record<string, string>;
	    npcs: string[];
	    is_entry_point: boolean;
	    display_name: string;
	    is_inter_cluster: boolean;
	    target_sub_cluster: string;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.sub_cluster_id = source["sub_cluster_id"];
	        this.type = source["type"];
	        this.connections = source["connections"];
	        this.connection_names = source["connection_names"];
	        this.npcs = source["npcs"];
	        this.is_entry_point = source["is_entry_point"];
	        this.display_name = source["display_name"];
	        this.is_inter_cluster = source["is_inter_cluster"];
	        this.target_sub_cluster = source["target_sub_cluster"];
	    }
	}
	export class Interaction {
	    id: string;
	    type: string;
	    content: string;
	    additional_content?: string;
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
	        this.additional_content = source["additional_content"];
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
	export class TemperTraits {
	    id: string;
	    npc_id: string;
	    prudence: number;
	    emotionality: number;
	    independence: number;
	    optimism: number;
	    flexibility: number;
	    aggressiveness: number;
	
	    static createFrom(source: any = {}) {
	        return new TemperTraits(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.npc_id = source["npc_id"];
	        this.prudence = source["prudence"];
	        this.emotionality = source["emotionality"];
	        this.independence = source["independence"];
	        this.optimism = source["optimism"];
	        this.flexibility = source["flexibility"];
	        this.aggressiveness = source["aggressiveness"];
	    }
	}
	export class NPC {
	    id: string;
	    name: string;
	    race: string;
	    location_id: string;
	    description: string;
	    temper_traits: TemperTraits;
	
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
	        this.temper_traits = this.convertValues(source["temper_traits"], TemperTraits);
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


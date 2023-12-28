import { expect } from 'chai'

import {
    configDirectory,
    organization,
} from '../src/config.js'

describe('config', () => { 
    it('configDirectory exitsts', () => { 
        expect(configDirectory).to.exist
    })
    it('organization exitsts', () => { 
        expect(organization).to.exist
    })
})

import { Tooltip } from 'bootstrap';
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { element } from 'tsx-vanilla';

import { Component } from './component';

export class SocialLinks extends Component {
	static buildGitHubLink(): Element {
		const anchor = (
			<a
				href="https://github.com/WoWLegacySims/wotlk"
				target="_blank"
				className="github-link link-alt"
				dataset={{ bsToggle: 'tooltip', bsTitle: 'Contribute on GitHub' }}>
				<i className="fab fa-github fa-lg" />
			</a>
		);
		Tooltip.getOrCreateInstance(anchor);
		return anchor;
	}
}
